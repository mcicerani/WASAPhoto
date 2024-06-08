<script>
import axios from 'axios'; // Importa l'istanza personalizzata di Axios



export default {
    data() {
        return {
            userStream: [] // Array to store the user stream photos
        };
    },
    mounted() {
        this.fetchUserStream(); // Fetch user stream photos when the component is mounted
    },
    methods: {
        async fetchUserStream() {
            try {
                const response = await axios.get('/api/user/stream'); // Make a GET request to the user stream API endpoint
                this.userStream = response.data; // Store the fetched user stream photos in the data property
            } catch (error) {
                console.error(error);
            }
        }
    }
};
</script>

<template>
    <div>
        <h1>User Stream</h1>
        <div v-if="userStream.length === 0">No photos in user stream</div>
        <div v-else>
            <div v-for="photo in userStream" :key="photo.id">
                <img :src="photo.url" alt="User Stream Photo">
            </div>
        </div>
    </div>
</template>