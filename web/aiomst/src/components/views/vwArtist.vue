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
        <eleCardList req="albums"></eleCardList>
      </v-row>
    </v-container>
    <v-container v-else fluid style="width: 75%;">
      <eleCardList req="artists" v-if="artists"></eleCardList>
    </v-container>
  </div>
</template>

<script>
import {mapGetters, mapState} from 'vuex';
import eleCardList  from '@/components/elements/eleCardList.vue';
import cmpntLoadBar from '@/components/layout/cmpntLoadBar.vue';

export default {
  name: 'Artist',
  components: {eleCardList, cmpntLoadBar},
  data: function() {
    return {
      finishedLoading: false,
      isSingleArtistView: false,
    }
  },
  computed: {
    ...mapGetters({
      artists:       'getArtists',
      currentArtist: 'currentArtist',
    }),
    ...mapState({
      apiResp: 'gResp', 
    })
  },
  watch: {
    apiResp: function() {
      this.isSingleArtistView = false
      if (this.currentArtist) {
        this.isSingleArtistView = true
      }
      this.finishedLoading = true
    }
  }
}
</script>