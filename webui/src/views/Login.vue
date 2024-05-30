<script>

import axios from "axios";

export default {
    data() {
        return {
            username: "",
            password: ""
        }
    },
    methods: {
        async doLogin() {
            try {
                const response = await axios.post("/session", {
                    username: this.username,
                    });
                console.log('Login successful', response.data);
                const userId = response.data.userId;
                this.$router.push('/users/${userId}')
            } catch (error) {
                console.error('Error logging in', error);
            }
        }
    }
}
</script>

<template>
    <div>
        <h2>Login</h2>
        <form @submit.prevent="login">
            <div>
                <label for="username">Username</label>
                <input type="text" id="username" placeholder="Username" v-model="username">
            </div>
            <button type="submit">Login</button>
        </form>
    </div>
</template>