<template>
  <div class="current-track">
    <!-- <div class="current-track_cover">

    </div>
    <div class="current-track_info">
      <router-link
        class="current-track_name"
        :to="{ name: 'album', params: { id: playback.item.album.id } }"
      > 
        {{ playback.item.name }}
      </router-link>

      <track-addition
        :trackID="currentTrackID"
        :isSaved="isSavedTrack"
        v-on:updateTrackstatus="onTrackUpdate"
      />

      <div class="current-track_artist">
        <router-link
          class="current-track_link"
          v-for="(artist, index) in playback.item.artists"
          :key="index"
          :to="{ name: 'artist', params: { id: artist.id } }"
        >
          {{ artist.name }}
          <template v-if="index !== playback.item.artists.length - 1">
            ,&nbsp;
          </template>
        </router-link>
      </div> -->
  </div>
</template>

<script>
export default {
  
}
</script>