const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

async function getJson(path) {
  const res = await fetch(`${API_BASE}${path}`)
  if (!res.ok) {
    throw new Error(`${path} responded with ${res.status}`)
  }
  return res.json()
}

export async function fetchProjects() {
  const data = await getJson('/api/projects')
  return data.projects
}

export async function fetchCertifications() {
  const data = await getJson('/api/certifications')
  return data.certifications
}

export async function fetchGithubStats() {
  return getJson('/api/github-stats')
}
