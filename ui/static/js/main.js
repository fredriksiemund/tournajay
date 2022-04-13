async function deleteTournament(id) {
  try {
    await fetch(`/tournament/${id}`, { method: "DELETE" });
    window.location.href = "/";
  } catch (err) {
    alert("Failed to delete tournament");
  }
}

async function deleteParticipant(tournamentId, userId) {
  try {
    await fetch(`/tournament/${tournamentId}/participants/${userId}`, {
      method: "DELETE",
    });
    window.location.reload();
  } catch (err) {
    alert("Failed to remove participant");
  }
}
