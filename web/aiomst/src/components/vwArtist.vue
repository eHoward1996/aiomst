<template>
  <cmpntLoadBar v-if="!finishedLoading"></cmpntLoadBar>
  <div v-else>
    <v-container v-if="isSingleArtistView" fluid style="width: 75%;">
      <v-row>
        <v-col>
          <h1>{{this.artists[0].title}} Albums</h1>
        </v-col>
      </v-row>
      <v-row>
        <cmpntCardList req="albums"></cmpntCardList>
      </v-row>
    </v-container>
    <v-container v-else fluid style="width: 75%;">
      <cmpntCardList req="artists" v-if="artists"></cmpntCardList>
    </v-container>
  </div>
</template>

<script>
import cmpntCardList from '@/components/cmpntCardList.vue';
import cmpntLoadBar from '@/components/cmpntLoadBar.vue';

export default {
  name: 'Artist',
  components: {cmpntCardList, cmpntLoadBar},
  created: function() {
    let navInfo = {
      path: '/artists',
      params: this.$route.query,
    };

    this.$store.dispatch('makeApiRequest', navInfo).then(() => {
      this.artists = this.$store.getters.artists;
      this.albums = this.$store.getters.albums;
      this.isSingleArtistView = false;
      if (this.artists.length === 1) {
        this.isSingleArtistView = true;
      }
      this.finishedLoading = true;
    });
    this.active_tab = 1;
  },
  data: function() {
    return {
      active_tab: 1,
      artists: null,
      albums: null,
      finishedLoading: false,
      isSingleArtistView: false,
    }
  },
  watch: {
    '$route.query.id': function() {
      this.finishedLoading = false;
      if (this.$route.query.id !== undefined) {
        this.artists = this.artists.filter(a => a.id === Number(this.$route.query.id));
        this.albums = this.albums.filter(a => a.artist_id === Number(this.$route.query.id));
        this.isSingleArtistView = true;
        this.finishedLoading = true;
        return
      } 
      
      this.$store.dispatch('makeApiRequest', {path: '/artists'}).then(() => {
        this.artists = this.$store.getters.artists;
        this.albums = this.$store.getters.albums;
        this.isSingleArtistView = false;
        this.finishedLoading = true;
      });
    }
  }
}
</script>