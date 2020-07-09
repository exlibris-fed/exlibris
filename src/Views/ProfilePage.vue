<template>
  <b-container>
    <b-spinner v-if="user == null" />
    <Profile
      v-else
      :user="user"
    />
  </b-container>
</template>

<script>
import Profile from '../components/Profile.vue'

export default {
  name: 'ProfilePage',

  components: {
    Profile
  },

  props: {
    axios: {
      type: Function,
      required: true
    },
    user: {
      type: Object,
      default: null
    }
  },

  data: function () {
    return {
      feed: null
    }
  }
}
</script>
