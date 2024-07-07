<script>
import api from "@/services/axios";

export default {
    data() {
        return {
            userStream: [] // Initialize userStream as an empty array
        };
    },
    mounted() {
        this.fetchUserStream(); // Fetch user stream photos when the component is mounted
    },
    methods: {
        async fetchUserStream() {
            try {
                const loggedInUserId = localStorage.getItem('loggedInUserId');
                const response = await api.get(`/users/${loggedInUserId}/stream`, {// Make a GET request to the user stream API endpoint
                    headers :{
                        Authorization : localStorage.getItem("token")
                    }
                });
                this.userStream = response.data.Photos; // Store the fetched user stream photos in the data property
            } catch (error) {
                console.error(error);
                // Handle error if needed
            }
        }
    }
};
</script>


<template>
    <div class="row mt-4">
        <div class="col-md-12">
            <ul class="listaFoto">
                <li v-for="(photo, index) in userStream" :key="photo.id" :class="{'new-row': index % 1 === 0}">
                    <div class="card">
                        <img class="text-center" :src="'data:image/jpeg;base64,' + photo.image_data" alt="User Photo">
                        <div class="text-center likes">
                            <p>
                                <svg class="feather">
                                    <use href="/feather-sprite-v4.29.0.svg#heart"/>
                                </svg>
                                {{ photo.likes }}
                            </p>
                        </div>
                        <div class="text-center comments">
                            <p>
                                <svg class="feather">
                                    <use href="/feather-sprite-v4.29.0.svg#message-square"/>
                                </svg>
                                {{ photo.comments }}
                            </p>
                        </div>
                    </div>
                </li>
            </ul>
        </div>
    </div>
</template>
