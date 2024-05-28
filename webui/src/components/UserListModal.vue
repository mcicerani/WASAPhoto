<template>
    <div class="user-list-modal">
        <h2>{{ title }}</h2>
        <ul>
            <li v-for="user in userList" :key="user.id">
                {{ user.username }}
            </li>
        </ul>
    </div>
</template>

<script>
export default {
    props: {
        title: {
            type: String,
            required: true
        },
        userList: {
            type: Array,
            required: true
        }
    }
}
</script>