package sqlreflect

import (
	"testing"

	"github.com/Masterminds/squirrel"
)

func loadTestTable(t *testing.T, name string) *Table {
	si := New(DBOptions{Driver: "postgres", Queryer: squirrel.NewStmtCacheProxy(db)})

	tt, err := si.Table(name, tCatalog, "")
	if err != nil {
		t.Fatal(err)
	}
	if tt.TableNameField != name {
		t.Fatalf("Expected %q , got %q", name, tt.TableNameField)
	}
	return tt
}

func TestTable_Privileges(t *testing.T) {
	table := loadTestTable(t, "person")
	// TODO: Actually call Privileges and check the output.
	privs, err := table.Privileges()
	if err != nil {
		t.Fatal(err)
	}
	if len(privs) == 0 {
		t.Fatalf("Expected at least one privilege on %s", table.TableNameField)
	}

	if !privs[0].IsGrantable.Bool {
		t.Errorf("Expected %v to be grantable.", privs[0])
	}
}

func TestTable_Constraints(t *testing.T) {
	table := loadTestTable(t, "person")
	// TODO: Actually call Privileges and check the output.
	constraints, err := table.Constraints()
	if err != nil {
		t.Fatal(err)
	}
	if len(constraints) == 0 {
		t.Fatalf("Expected at least one privilege on %s", table.TableNameField)
	}
	for _, c := range constraints {
		t.Logf("CONSTRAINT: %v", c)
	}
}

func TestTable_ConstraintsByType(t *testing.T) {
	table := loadTestTable(t, "person")
	constraints, err := table.ConstraintsByType(ConstraintPrimaryKey)
	if err != nil {
		t.Fatal(err)
	}
	if len(constraints) != 1 {
		t.Fatalf("Expected one primary key on %s", table.TableNameField)
	}
	constraints, err = table.ConstraintsByType(ConstraintForeignKey)
	if err != nil {
		t.Fatal(err)
	}
	if len(constraints) != 0 {
		t.Fatalf("Expected 0 constraints, got %d", len(constraints))
	}
}

func TestTable_Consraint(t *testing.T) {
	table := loadTestTable(t, "org")
	con, err := table.Constraint("org_pkey")
	if err != nil {
		t.Fatal(err)
	}

	if con.TableNameField != "org" {
		t.Errorf("Expected table org, got %q", con.TableNameField)
	}

	_, err = table.Constraint("org_no_such_key")
	if err == nil {
		t.Error("Expected org_no_such_key lookup to produce an error.")
	}
}

func TestTable_PrimaryKey(t *testing.T) {
	table := loadTestTable(t, "person")
	pk, err := table.PrimaryKey()
	if err != nil {
		t.Fatal(err)
	}
	if pk.ConstraintType != ConstraintPrimaryKey {
		t.Fatalf("Unexpected primary key constraint type: %q", pk.ConstraintType)
	}
}

func TestTable_ForeignKeys(t *testing.T) {
	table := loadTestTable(t, "person")
	fk, err := table.ForeignKeys()
	if err != nil {
		t.Fatal(err)
	}
	if len(fk) != 0 {
		t.Errorf("Expected no foreign keys for the person table, but got %d", len(fk))
	}

	table = loadTestTable(t, "employees")
	fk, err = table.ForeignKeys()
	if err != nil {
		t.Fatal(err)
	}
	if len(fk) != 2 {
		t.Fatalf("Expected 2 foreign keys for the employees table, but got %d", len(fk))
	}
}

func TestTable_InViews(t *testing.T) {
	table := loadTestTable(t, "person")
	views, err := table.InViews()
	if err != nil {
		t.Fatal(err)
	}
	if len(views) != 1 {
		t.Errorf("Expected 1 view, got %d", len(views))
	}

	if views[0].ViewDefinition == "" {
		t.Errorf("View was not initialized")
	}
}

func TestTable_Columns(t *testing.T) {
	table := loadTestTable(t, "person")
	cols, err := table.Columns()
	if err != nil {
		t.Fatal(err)
	}
	if l := len(cols); l != 3 {
		t.Errorf("Expected 3 columns, got %d", l)
	}

	for _, c := range cols {
		switch c.OrdinalPosition {
		case 1:
			if c.Name != "id" {
				t.Errorf("Expected ordinal position 1 to be id, got %q", c.Name)
			}
			if c.IsNullable.Bool {
				t.Error("Expected id to not be nullable")
			}
			if c.DataType != "integer" {
				t.Errorf("Expected id DataType to be integer, got %q", c.DataType)
			}
		case 2:
			if c.Name != "first_name" {
				t.Errorf("Expected ordinal position 2 to be first_name, got %q", c.Name)
			}
			if !c.IsNullable.Bool {
				t.Error("Expected id to be nullable")
			}
			if c.DataType != "character varying" {
				t.Errorf("Expected first_name DataType to be varchar, got %q", c.DataType)
			}
		}
	}
}

func TestTable_Column(t *testing.T) {
	table := loadTestTable(t, "person")
	c, err := table.Column("id")
	if err != nil {
		t.Fatal(err)
	}

	if c.Name != "id" {
		t.Errorf("Expected ordinal position 1 to be id, got %q", c.Name)
	}
	if c.IsNullable.Bool {
		t.Error("Expected id to not be nullable")
	}
	if c.DataType != "integer" {
		t.Errorf("Expected id DataType to be integer, got %q", c.DataType)
	}
}
