<script>

import axios from 'axios'; // Importa l'istanza personalizzata di Axios

// inserisce nuovo username e lo sostiutisce a quello vecchio


</script>

<template>
    <form>
        <div>
            <label for="username">Username</label>
            <input type="text" id="username" v-model="user.username">
        </div>
        <button type="submit">Salva</button>
    </form>
</template>