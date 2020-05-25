<template>
  <v-content>
    <v-container fluid>
    </v-container>
  </v-content>
</template>

<script>
// import rndArtistsList from '@/components/rndArtistsList.vue';
// import rndAlbumsList  from '@/components/rndAlbumsList.vue';
// import rndSongsList   from '@/components/rndSongsList.vue';


export default {
  name: 'Content',
  components: {
    // rndArtistsList,
    // rndAlbumsList, 
    // rndSongsList
  },
  props: ['page'],
  data: () => ({
    // respTypes: ['artists', 'albums', 'songs'],
  }),
  created() {
    this.$vuetify.theme.dark = true
  },
}
</script>