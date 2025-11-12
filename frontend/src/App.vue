<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
    <nav class="bg-white shadow-md">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center gap-2">
            <BookOpenIcon class="w-8 h-8 text-blue-600" />
            <h1 class="text-2xl font-bold text-gray-900">VicNotes</h1>
          </div>
          <div class="flex items-center gap-4">
            <router-link 
              v-if="!authStore.isAuthenticated"
              to="/login" 
              class="text-gray-600 hover:text-gray-900"
            >
              Login
            </router-link>
            <router-link 
              v-if="!authStore.isAuthenticated"
              to="/register" 
              class="btn-primary"
            >
              Sign Up
            </router-link>
            <div v-if="authStore.isAuthenticated" class="flex items-center gap-4">
              <span class="text-gray-600">{{ authStore.user?.email }}</span>
              <button @click="logout" class="btn-secondary">
                Logout
              </button>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <router-view />
    </main>
  </div>
</template>

<script>
import { defineComponent } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { BookOpen } from 'lucide-vue-next'

export default defineComponent({
  name: 'App',
  components: {
    BookOpenIcon: BookOpen
  },
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()

    const logout = () => {
      authStore.logout()
      router.push('/login')
    }

    return {
      authStore,
      logout
    }
  }
})
</script>
