<template>
  <b-nav-item-dropdown
    right
  >
    <template
      v-slot:button-content
    >
      <UserIcon
        :user="user"
        height="50"
        round
      />
      {{ user.name }}
    </template>

    <b-dropdown-item
      :to="profileLink"
    >
      {{ $t('profile') }}
    </b-dropdown-item>
    <b-dropdown-item
      to="logout"
    >
      {{ $t('logout') }}
    </b-dropdown-item>
  </b-nav-item-dropdown>
</template>

<script>
import UserIcon from './UserIcon.vue'

export default {
  name: 'HeaderProfile',
  components: {
    UserIcon
  },
  props: {
    user: {
      type: Object,
      required: true
    }
  },
  computed: {
    profileLink () {
      return '/@' + this.user.preferredUsername
    }
  },
  i18n: {
    messages: {
      en: {
        profile: 'Profile',
        logout: 'Log Out'
      }
    }
  }
}
</script>
