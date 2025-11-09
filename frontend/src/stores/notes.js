import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api/client'

export const useNotesStore = defineStore('notes', () => {
  const notes = ref([])
  const loading = ref(false)
  const error = ref(null)

  const fetchNotes = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await api.get('/api/v1/notes')
      notes.value = response.data || []
    } catch (err) {
      error.value = err.response?.data?.message || 'Failed to fetch notes'
      console.error(error.value)
    } finally {
      loading.value = false
    }
  }

  const createNote = async (title, content) => {
    try {
      const response = await api.post('/api/v1/notes', { title, content })
      notes.value.unshift(response.data)
      return response.data
    } catch (err) {
      error.value = err.response?.data?.message || 'Failed to create note'
      throw err
    }
  }

  const updateNote = async (id, title, content) => {
    try {
      await api.put(`/api/v1/notes/${id}`, { title, content })
      const index = notes.value.findIndex(n => n.id === id)
      if (index !== -1) {
        notes.value[index] = { ...notes.value[index], title, content }
      }
    } catch (err) {
      error.value = err.response?.data?.message || 'Failed to update note'
      throw err
    }
  }

  const deleteNote = async (id) => {
    try {
      await api.delete(`/api/v1/notes/${id}`)
      notes.value = notes.value.filter(n => n.id !== id)
    } catch (err) {
      error.value = err.response?.data?.message || 'Failed to delete note'
      throw err
    }
  }

  return {
    notes,
    loading,
    error,
    fetchNotes,
    createNote,
    updateNote,
    deleteNote
  }
})
