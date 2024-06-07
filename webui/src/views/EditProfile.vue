<template>
    <div class="modal">
        <form>
            <label for="new-username">Username:</label>
            <input type="text" id="new-username" placeholder="Username" v-model="username">
            <button @click="setMyUsername">Update</button>
        </form>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    data() {
        return {
            username: ''
        }
    },
    methods: {
        setMyUsername() {
            axios.post('/users/{userid}', {
                username: this.username
            }).then(() => {
                this.$emit('username-updated', this.username);
            });
        }
    }
}
</script>