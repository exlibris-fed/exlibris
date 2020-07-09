<template>
  <div>
    <UserIcon
      :user="user"
    />
    <h1>
      {{ user.name }}
      <small class="text-muted">{{ $t('at') }}{{ user.preferredUsername }}</small>
    </h1>

    <div
      v-if="feed"
    >
      <h2>{{ $t('feed') }}</h2>
      <BookList :feed="feed" />
    </div>
  </div>
</template>

<script>
import BookList from './BookList.vue'
import UserIcon from './UserIcon.vue'

export default {
  name: 'Profile',
  components: {
    BookList,
    UserIcon
  },
  props: {
    user: {
      type: Object,
      required: true
    },
    feed: {
      type: Array,
      default: () => []
    }
  },
  i18n: {
    messages: {
      en: {
        at: '@',
        feed: 'Recently Read'
      }
    }
  }
}
</script>
