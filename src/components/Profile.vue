<template>
  <div>
    <h1>
      {{ user.name }}
      <small class="text-muted">{{ user.preferredUsername }}</small>
    </h1>
  </div>
</template>

<script>
export default {
  name: 'Profile',
  props: {
    user: {
      type: Object,
      required: true
    }
  }
}
</script>
