<template>
  <b-col
    cols="12"
    md="3"
  >
    <b-card
      :title="book.title"
      :img-src="coverimage"
    >
      <b-card-text>
        {{ $t('attribution') }} {{ book.authors[0] }}
      </b-card-text>
      <b-button
        @click="readBook"
      >
        {{ $t('read') }}
      </b-button>
    </b-card>
  </b-col>
</template>

<script>
export default {
  name: 'BookGridCard',
  i18n: {
    messages: {
      en: {
        attribution: 'By',
        read: "I've Read This",
        review: 'Write A Review'
      }
    }
  },
  props: {
    book: {
      type: Object,
      required: true
    }
  },
  computed: {
    coverimage: function () {
      return this.book.covers.large || this.book.covers.medium || this.book.covers.small
    }
  },
  methods: {
    readBook: function () {
      this.$emit('read', this.book)
    }
  }
}
</script>
