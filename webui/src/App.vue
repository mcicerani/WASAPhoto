<script setup>
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { ref, onMounted, watch } from 'vue'

const loggedInUserId = ref(null);
const router = useRouter();

const token = localStorage.getItem('token');

if (token) {
  onMounted(() => {
    const loggedInUserId = loadUserIdFromToken();
    router.push(`/users/${loggedInUserId}/profile`); // Usa il template literal corretto
  });

  watch(router.currentRoute, () => {
    loadUserIdFromToken();
  });
} else {
  console.log("Token not found in localStorage");
}

function loadUserIdFromToken() {
  const token = localStorage.getItem('token');
  if (token) {
    loggedInUserId.value = token.split(' ')[1]; // Estrai l'ID dell'utente dal token
    console.log("User ID extracted from token:", loggedInUserId.value);
  } else {
    console.log("Token not found in localStorage");
  }
}

function logout() {
  localStorage.removeItem('token');
  localStorage.removeItem('username');
  localStorage.removeItem('loggedInUserId');
  console.log("User logged out");
  router.push('/session');
}

console.log("Component setup complete, userId:", loggedInUserId.value);
</script>

<template>
  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASAPhoto</a>
    <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
    <button class="logout" @click="logout">Logout</button>
  </header>

  <div class="container-fluid">
    <div class="row">
      <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
        <div class="position-sticky pt-3 sidebar-sticky">
          <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
          </h6>
          <ul class="nav flex-column">
            <li class="nav-item">
              <RouterLink :to="`/users/${loggedInUserId}/profile`" class="nav-link">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
                Profilo
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink :to="`/users/${loggedInUserId}/stream`" class="nav-link">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
                Stream
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink to="/users" class="nav-link">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
                Ricerca Utenti
              </RouterLink>
            </li>
          </ul>
        </div>
      </nav>

      <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<style>
/* Aggiungi i tuoi stili personalizzati qui */
</style>
