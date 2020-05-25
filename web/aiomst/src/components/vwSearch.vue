<template>
  <v-container>
    <v-row v-if="$store.state.gResp">
      <rndArtistsList></rndArtistsList>
    </v-row>
    <v-row v-if="$store.state.gResp">
      <rndAlbumsList></rndAlbumsList>
    </v-row>
    <v-row v-if="$store.state.gResp">
      <rndSongsList></rndSongsList>
    </v-row>
  </v-container>
</template>

<script>
import rndArtistsList from '@/components/rndArtistsList.vue';
import rndAlbumsList from '@/components/rndAlbumsList.vue';
import rndSongsList from '@/components/rndSongsList.vue';

export default {
  name: 'Search',
  components: {rndArtistsList, rndAlbumsList, rndSongsList},
}
</script>