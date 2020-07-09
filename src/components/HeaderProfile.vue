<template>
  <b-nav-item-dropdown
    right
  >
    <template
      v-slot:button-content
    >
      <b-img
        width="50"
        height="50"
        alt=""
        :src="icon"
        rounded="circle"
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
export default {
  name: 'HeaderProfile',
  props: {
    user: {
      type: Object,
      required: true
    }
  },
  computed: {
    icon () {
      return (this.user && this.user.icon) || require('../../public/default-profile.jpg')
    },
    profileLink () {
      return '@' + this.user.preferredUsername
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
