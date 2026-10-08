package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"siakad-mini/app/model"
)

type CourseFilter struct {
	Semester  *int
	Search    string
	Available bool
}

type CourseRepository interface {
	List(ctx context.Context, f CourseFilter) ([]model.Course, error)
	GetForEnrollment(ctx context.Context, tx pgx.Tx, id int64) (model.Course, error)
}

type coursePostgresRepository struct{ pool *pgxpool.Pool }

func NewCourseRepository(pool *pgxpool.Pool) CourseRepository {
	return &coursePostgresRepository{pool: pool}
}

func (r *coursePostgresRepository) List(ctx context.Context, f CourseFilter) ([]model.Course, error) {
	where := []string{"1=1"}
	args := []any{}
	n := 1
	if f.Semester != nil {
		where = append(where, fmt.Sprintf("c.semester=$%d", n))
		args = append(args, *f.Semester)
		n++
	}
	if f.Search != "" {
		where = append(where, fmt.Sprintf("(c.kode_mk ILIKE $%d OR c.nama_mk ILIKE $%d)", n, n))
		args = append(args, "%"+f.Search+"%")
		n++
	}
	q := `SELECT c.id,c.kode_mk,c.nama_mk,c.sks,c.semester,c.kuota,COUNT(e.id)::int AS terisi,
        (c.kuota-COUNT(e.id))::int AS sisa_kuota
        FROM courses c LEFT JOIN enrollments e ON e.course_id=c.id
        WHERE ` + strings.Join(where, " AND ") + ` GROUP BY c.id`
	if f.Available {
		q += ` HAVING (c.kuota-COUNT(e.id)) > 0`
	}
	q += ` ORDER BY c.semester,c.kode_mk`
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Course{}
	for rows.Next() {
		var c model.Course
		if err := rows.Scan(&c.ID, &c.KodeMK, &c.NamaMK, &c.SKS, &c.Semester, &c.Kuota, &c.Terisi, &c.SisaKuota); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *coursePostgresRepository) GetForEnrollment(ctx context.Context, tx pgx.Tx, id int64) (model.Course, error) {
	var c model.Course
	err := tx.QueryRow(ctx, `SELECT id,kode_mk,nama_mk,sks,semester,kuota FROM courses WHERE id=$1 FOR UPDATE`, id).
		Scan(&c.ID, &c.KodeMK, &c.NamaMK, &c.SKS, &c.Semester, &c.Kuota)
	if err != nil {
		return c, err
	}
	if err := tx.QueryRow(ctx, `SELECT COUNT(*)::int FROM enrollments WHERE course_id=$1`, id).Scan(&c.Terisi); err != nil {
		return c, err
	}
	c.SisaKuota = c.Kuota - c.Terisi
	return c, nil
}
