package my_tests

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/entities"
	"darbelis.eu/stabas/util"
	"fmt"
	"reflect"
	"testing"
)

func TestAddParticipant(t *testing.T) {
	rep := &dao.ParticipantsRepository{}

	rep.AddParticipant(&entities.Participant{Name: "first"})
	rep.AddParticipant(&entities.Participant{Name: "second"})
	rep.AddParticipant(&entities.Participant{Name: "third"})

	got := rep.GetParticipants()
	want := []*entities.Participant{
		{Id: 1, Name: "first"},
		{Id: 2, Name: "second"},
		{Id: 3, Name: "third"},
	}

	assertParticipantsArrays(t, want, got)
}
func TestUpdateParticipant(t *testing.T) {
	rep := &dao.ParticipantsRepository{}

	rep.AddParticipant(&entities.Participant{Name: "first"})
	rep.AddParticipant(&entities.Participant{Name: "second"})
	rep.AddParticipant(&entities.Participant{Name: "third"})
	rep.UpdateParticipant(&entities.Participant{Id: 1, Name: "first modified"})
	rep.RemoveParticipant(2)

	got := rep.GetParticipants()
	want := []*entities.Participant{
		{Id: 1, Name: "first modified"},
		{Id: 3, Name: "third"},
	}

	assertParticipantsArrays(t, want, got)
}

func assertParticipantsArrays(t *testing.T, got, want []*entities.Participant) {
	if !reflect.DeepEqual(want, got) {
		gotRepresentations := util.ArrayMap(got, func(p *entities.Participant) string { return fmt.Sprintf("%d:%s", p.Id, p.Name) })
		wantRepresentations := util.ArrayMap(want, func(p *entities.Participant) string { return fmt.Sprintf("%d:%s", p.Id, p.Name) })

		t.Errorf("Arrays were not equal want %q,  got %q", wantRepresentations, gotRepresentations)
	}
}
