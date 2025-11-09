<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-4xl font-bold">My Notes</h1>
      <button @click="showCreateModal = true" class="btn-primary flex items-center gap-2">
        <Plus class="w-5 h-5" />
        New Note
      </button>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="card w-full max-w-2xl max-h-96 overflow-y-auto">
        <h2 class="text-2xl font-bold mb-4">{{ editingNote ? 'Edit Note' : 'Create New Note' }}</h2>
        
        <form @submit.prevent="saveNote" class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-2">Title</label>
            <input
              v-model="noteForm.title"
              type="text"
              class="input-field"
              placeholder="Note title"
              required
            />
          </div>

          <div>
            <label class="block text-sm font-medium mb-2">Content</label>
            <textarea
              v-model="noteForm.content"
              class="input-field h-40 resize-none"
              placeholder="Write your note here..."
            ></textarea>
          </div>

          <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
            {{ error }}
          </div>

          <div class="flex gap-3 justify-end">
            <button
              type="button"
              @click="closeModal"
              class="btn-secondary"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="notesStore.loading"
              class="btn-primary disabled:opacity-50"
            >
              {{ notesStore.loading ? 'Saving...' : 'Save Note' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="notesStore.loading && notesStore.notes.length === 0" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      <p class="mt-4 text-gray-600">Loading notes...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="notesStore.notes.length === 0" class="text-center py-12">
      <BookOpen class="w-16 h-16 text-gray-300 mx-auto mb-4" />
      <p class="text-gray-600 text-lg">No notes yet. Create your first note!</p>
    </div>

    <!-- Notes Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="note in notesStore.notes"
        :key="note.id"
        class="card group cursor-pointer hover:shadow-xl"
      >
        <div class="flex justify-between items-start mb-3">
          <h3 class="text-lg font-semibold flex-1 line-clamp-2">{{ note.title }}</h3>
          <div class="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              @click="editNote(note)"
              class="p-2 hover:bg-gray-100 rounded-lg"
              title="Edit"
            >
              <Edit2 class="w-4 h-4 text-blue-600" />
            </button>
            <button
              @click="deleteNoteConfirm(note.id)"
              class="p-2 hover:bg-gray-100 rounded-lg"
              title="Delete"
            >
              <Trash2 class="w-4 h-4 text-red-600" />
            </button>
          </div>
        </div>

        <p class="text-gray-600 line-clamp-4 mb-4">{{ note.content || 'No content' }}</p>

        <p class="text-xs text-gray-400">
          {{ formatDate(note.created_at) }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useNotesStore } from '../stores/notes'
import { Plus, Edit2, Trash2, BookOpen } from 'lucide-vue-next'

const notesStore = useNotesStore()
const showCreateModal = ref(false)
const editingNote = ref(null)
const error = ref(null)

const noteForm = ref({
  title: '',
  content: ''
})

onMounted(() => {
  notesStore.fetchNotes()
})

const closeModal = () => {
  showCreateModal.value = false
  editingNote.value = null
  noteForm.value = { title: '', content: '' }
  error.value = null
}

const editNote = (note) => {
  editingNote.value = note
  noteForm.value = {
    title: note.title,
    content: note.content
  }
  showCreateModal.value = true
}

const saveNote = async () => {
  error.value = null

  if (!noteForm.value.title.trim()) {
    error.value = 'Title is required'
    return
  }

  try {
    if (editingNote.value) {
      await notesStore.updateNote(
        editingNote.value.id,
        noteForm.value.title,
        noteForm.value.content
      )
    } else {
      await notesStore.createNote(
        noteForm.value.title,
        noteForm.value.content
      )
    }
    closeModal()
  } catch (err) {
    error.value = err.response?.data?.message || 'Failed to save note'
  }
}

const deleteNoteConfirm = async (id) => {
  if (confirm('Are you sure you want to delete this note?')) {
    try {
      await notesStore.deleteNote(id)
    } catch (err) {
      alert(err.response?.data?.message || 'Failed to delete note')
    }
  }
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  })
}
</script>
