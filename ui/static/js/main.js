async function deleteTournament(id) {
  try {
    await fetch(`/tournament/${id}`, { method: "DELETE" });
    window.location.href = "/";
  } catch (err) {
    alert("Failed to delete tournament");
  }
}
