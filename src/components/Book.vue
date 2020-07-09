<template>
  <b-col
    v-if="type === 'card'"
    cols="12"
    md="3"
    class="mb-3"
  >
    <b-card
      :title="book.title"
    >
      <b-link
        :to="{ name: 'book', params: { book: idSlug } }"
      >
        <b-card-img-lazy
          :src="coverImageLg"
        />
      </b-link>
      <b-card-text>
        <p
          v-if="published"
        >
          {{ $t('published') }} {{ published }}
        </p>
        <p
          v-if="authors"
        >
          {{ $t('attribution') }} {{ authors }}
        </p>
      </b-card-text>
      <b-button
        @click="readBook"
      >
        {{ $t('read') }}
      </b-button>
    </b-card>
  </b-col>

  <b-list-group-item
    v-else-if="type === 'list'"
    :to="{ name: 'book', params: { book: idSlug } }"
  >
    <h3>{{ book.title }}</h3>
    <p
      v-if="authors"
      class="text-muted"
    >
      {{ $t('attribution') }} {{ authors }}
    </p>
  </b-list-group-item>

  <div v-else>
    <b-img
      :src="coverImageLg"
      left
    />
    <h1> {{ book.title }}</h1>
    <p
      v-if="authors"
      class="text-muted"
    >
      {{ $t('attribution') }} {{ authors }}
    </p>
    <p
      v-if="published"
    >
      {{ $t('published') }} {{ published }}
    </p>
    <b-button
      @click="readBook"
    >
      {{ $t('read') }}
    </b-button>

    <div class="clearfix" />

    <blockquote class="blockquote">
      {{ book.description }}
    </blockquote>
  </div>
</template>

<script>
export default {
  name: 'Book',
  i18n: {
    messages: {
      en: {
        attribution: 'By',
        published: 'Published',
        read: "I've Read This",
        review: 'Write A Review'
      }
    }
  },
  props: {
    book: {
      type: Object,
      required: true
    },
    type: {
      type: String,
      default: 'full'
    }
  },
  computed: {
    authors () {
      return this.book && this.book.authors && this.book.authors.join(', ')
    },

    coverImageLg: function () {
      return (this.book.covers && (this.book.covers.large || this.book.covers.medium || this.book.covers.small)) || require('../../public/default-cover.jpg')
    },

    coverImageSm: function () {
      return (this.book.covers && (this.book.covers.small || this.book.covers.medium || this.book.covers.large)) || require('../../public/default-cover.jpg')
    },

    idSlug () {
      return this.book.id.split('/').pop()
    },

    published () {
      // don't have a published date
      if (!this.book.published) {
        return null
      }
      // published date just so happent to be the unix epoch? fishy.
      const timestamp = new Date(this.book.published)
      if (timestamp.getTime() === 0) {
        return null
      }
      return timestamp.getFullYear()
    }
  },
  methods: {
    readBook: function () {
      this.$emit('read', this.book)
    }
  }
}
</script>
