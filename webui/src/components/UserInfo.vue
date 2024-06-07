<script>
import axios from 'axios';
import UserListModal from './UserListModal.vue';

export default {
    components: {
        UserListModal,
    },
    data() {
        return {
            username: '',
            numFollowers: 0,
            numFollows: 0
        }
    },
    methods: {
        followersModal() {
            axios.get('/users/{userid}/followers').then(response => {
                this.$modal.show(UserListModal, {
                    title: 'Followers',
                    userList: response.data
                }, {
                    height: 'auto',
                    width: 'auto'
                });
            });
        },
        FollowsModal() {
            axios.get('/users/{userid}/follows').then(response => {
                this.$modal.show(UserListModal, {
                    title: 'Follows',
                    userList: response.data
                }, {
                    height: 'auto',
                    width: 'auto'
                });
            });
        }
    },
    mounted() {
        axios.get('/users/{userid}').then(response => {
            this.username = response.data.username;
            this.numFollowers = response.data.numFollowers;
            this.numFollows = response.data.numFollows;
        });
    }

}
</script>

<template>
    <div>
        <h1>{{username}}</h1>
        <button>Edit Username</button>
        <button @click="followersModal">Followers: {{numFollowers}}</button>
        <button @click="FollowsModal">Follows: {{numFollows}}</button>
    </div>
</template>
