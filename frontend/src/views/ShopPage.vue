<template>
  <div class="page">
    <GlobalNav showBack @openSettings="showSettings = true" @openFriends="showFriends = true">
      <template #left>
        <button class="btn-back" @click="$router.push('/lobby')">← 返回大厅</button>
      </template>
    </GlobalNav>

    <main class="content">
      <ArtPlaceholder label="商城背景" width="100%" height="100%" bgColor="#0f0f1e" class="bg" />
      
      <div class="shop-container">
        <!-- 标签切换 -->
        <div class="tabs">
          <button :class="['tab', { active: shopTab === 'hero' }]" @click="shopTab = 'hero'">英雄商店</button>
          <button :class="['tab', { active: shopTab === 'skin' }]" @click="shopTab = 'skin'">皮肤商店</button>
        </div>

        <!-- 商品网格 -->
        <div class="shop-grid" v-if="shopTab === 'hero'">
          <div v-for="item in heroItems" :key="item.id" class="item-card">
            <ArtPlaceholder :label="item.name" width="100%" height="120" bgColor="#1a1a3e" />
            <div class="item-info">
              <span class="item-name">{{ item.name }}</span>
              <span class="item-price">💰 {{ item.price }}</span>
              <button v-if="item.owned" class="btn-owned" disabled>已拥有</button>
              <button v-else-if="userGold >= item.price" class="btn-buy" @click="buyItem(item)">购买</button>
              <button v-else class="btn-nocoin" disabled>金币不足</button>
            </div>
          </div>
        </div>

        <div class="shop-grid" v-if="shopTab === 'skin'">
          <div v-for="item in skinItems" :key="item.id" class="item-card">
            <ArtPlaceholder :label="item.name" width="100%" height="120" bgColor="#1a1a3e" />
            <div class="item-info">
              <span class="item-name">{{ item.name }}</span>
              <span class="item-belong">英雄: {{ item.heroName }}</span>
              <span class="item-price">💰 {{ item.price }}</span>
              <button v-if="item.owned" class="btn-owned" disabled>已拥有</button>
              <button v-else-if="userGold >= item.price" class="btn-buy" @click="buyItem(item)">购买</button>
              <button v-else class="btn-nocoin" disabled>金币不足</button>
            </div>
          </div>
        </div>
      </div>
    </main>

    <SettingsDialog :visible="showSettings" @close="showSettings = false" />
    <FriendsDialog :visible="showFriends" @close="showFriends = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/store/user'
import api from '@/api'
import GlobalNav from '@/components/GlobalNav.vue'
import ArtPlaceholder from '@/components/ArtPlaceholder.vue'
import SettingsDialog from '@/components/SettingsDialog.vue'
import FriendsDialog from '@/components/FriendsDialog.vue'

const userStore = useUserStore()
const userGold = computed(() => userStore.gold || 0)
const shopTab = ref<'hero' | 'skin'>('hero')
const showSettings = ref(false)
const showFriends = ref(false)

const heroItems = ref<any[]>([])
const skinItems = ref<any[]>([])

onMounted(async () => {
  try {
    const joinRoomInfo = await api.get('/shop/items')
    const allItems = ((joinRoomInfo as any).data || joinRoomInfo) || []
    heroItems.value = allItems.filter((item: any) => item.item_type === 'character').map((item: any) => ({
      id: item.id, shopItemId: item.id, name: item.name, price: item.price, owned: false,
    }))
    skinItems.value = allItems.filter((item: any) => item.item_type === 'skin').map((item: any) => ({
      id: item.id, shopItemId: item.id, name: item.name, heroName: '', price: item.price, owned: false,
    }))
    // 查询已拥有的
    try {
      const myRes: any = await api.get('/profile/characters')
      const myHeroes = (myRes.data || myRes) || []
      const ownedHeroIds = new Set(myHeroes.map((h: any) => h.id || h.character_id))
      heroItems.value.forEach((h: any) => { h.owned = ownedHeroIds.has(h.id) })
    } catch { /* ignore */ }
    try {
      const mySkinsRes: any = await api.get('/skins/my')
      const mySkins = (mySkinsRes.data || mySkinsRes) || []
      const ownedSkinIds = new Set(mySkins.map((s: any) => s.id || s.skin_id))
      skinItems.value.forEach((s: any) => { s.owned = ownedSkinIds.has(s.id) })
    } catch { /* ignore */ }
  } catch {
    heroItems.value = [
      { id: 1, name: '火焰战士', price: 1000, owned: false },
      { id: 2, name: '寒冰法师', price: 1200, owned: false },
    ]
    skinItems.value = [
      { id: 1, name: '暗夜皮肤', heroName: '火焰战士', price: 500, owned: false },
    ]
  }
})

async function buyItem(item: any) {
  try {
    await api.post('/shop/purchase', { item_id: item.shopItemId || item.id })
    item.owned = true
    await userStore.fetchUserInfo()
    alert('购买成功!')
  } catch (err: any) {
    alert(err?.data?.message || '购买失败')
  }
}
</script>

<style scoped>
.page { width: 100vw; height: 100vh; display: flex; flex-direction: column; overflow: hidden; }
.content { flex: 1; position: relative; padding: 20px; }
.bg { position: absolute; inset: 0; z-index: 0; }
.shop-container { position: relative; z-index: 1; }
.tabs { display: flex; gap: 8px; margin-bottom: 20px; }
.tab {
  padding: 10px 24px; background: rgba(255,255,255,.05); border: 1px solid #3a3a5e;
  border-radius: 8px; color: #888; cursor: pointer; font-size: 15px;
}
.tab.active { background: linear-gradient(135deg,#ffd70033,#ff8c0033); border-color: #ffd700; color: #ffd700; font-weight: bold; }
.shop-grid {
  display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px; max-height: calc(100vh - 200px); overflow-y: auto;
}
.item-card {
  background: rgba(20,20,40,.9); border: 1px solid #3a3a5e; border-radius: 12px;
  overflow: hidden; display: flex; flex-direction: column;
}
.item-info { padding: 12px; display: flex; flex-direction: column; gap: 6px; }
.item-name { color: #e0e0e0; font-size: 14px; font-weight: bold; }
.item-price { color: #ffd700; font-size: 14px; }
.item-belong { color: #888; font-size: 11px; }
.btn-buy {
  padding: 8px; border: none; border-radius: 6px; cursor: pointer;
  background: linear-gradient(135deg,#ffd700,#ff8c00); color: #1a1a1a; font-weight: bold;
}
.btn-nocoin { padding: 8px; border: none; border-radius: 6px; background: #2a2a4a; color: #666; }
.btn-owned { padding: 8px; border: none; border-radius: 6px; background: #22c55e44; color: #22c55e; }
.btn-back { background: rgba(255,255,255,.1); border: 1px solid #3a3a5e; color: #ccc; padding: 6px 16px; border-radius: 6px; cursor: pointer; }
</style>
