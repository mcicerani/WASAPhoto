<script setup>
import { ref, watch, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

const userId = ref(localStorage.getItem('loggedInUserId'));
const token = ref(localStorage.getItem('token'));
const username = ref(localStorage.getItem('username'));

const router = useRouter();

const isLoggedIn = computed(() => !!token.value);

function handleLoginSuccess(payload) {
  console.log('handleLoginSuccess called with payload:', payload); // Debugging

  username.value = payload.username;
  userId.value = payload.userId;
  token.value = payload.token;

  localStorage.setItem('username', payload.username);
  localStorage.setItem('loggedInUserId', payload.userId);
  localStorage.setItem('token', payload.token);

  console.log('Navigating to profile:', `/users/${userId.value}/profile`);
  router.push(`/users/${userId.value}/profile`);
}

function logout() {
  username.value = '';
  userId.value = '';
  token.value = '';
  localStorage.clear();
  router.push('/session');
}

watch(token, (newToken) => {
  if (newToken) {
    localStorage.setItem('token', newToken);
  } else {
    localStorage.removeItem('token');
  }
});

watch(userId, (newUserId) => {
  console.log('userId updated to:', newUserId); // Debugging
});

onMounted(() => {
  if (isLoggedIn.value && userId.value) {
    console.log('User already logged in, redirecting to profile:', `/users/${userId.value}/profile`);
    router.push(`/users/${userId.value}/profile`);
  } else {
    console.log('User not logged in, redirecting to session');
    router.push('/session');
  }
});
</script>

<template>
  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASAPhoto</a>
    <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
    <button v-if="isLoggedIn" class="logout" @click="logout">Logout</button>
  </header>

  <div class="container-fluid">
    <div class="row">
      <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
        <div class="position-sticky pt-3 sidebar-sticky">
          <ul v-if="isLoggedIn" class="nav flex-column">
            <li class="nav-item">
              <RouterLink :to="`/users/${userId}/profile`" class="nav-link">
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#home"/>
                </svg>
                Profilo
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink :to="`/users/${userId}/stream`" class="nav-link">
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#layout"/>
                </svg>
                Stream
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink to="/users" class="nav-link">
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#key"/>
                </svg>
                Ricerca Utenti
              </RouterLink>
            </li>
          </ul>
        </div>
      </nav>

      <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
        <router-view @login-success="handleLoginSuccess" />
      </main>
    </div>
  </div>
</template>

<style>
/* Add your custom styles here */
</style>
