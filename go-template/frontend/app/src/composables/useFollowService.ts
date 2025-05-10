import apiService from "@/services/ApiService";
import { ref } from "vue";
import { router } from "@/plugins/1.router";
import { ErrorPopup } from "@/utils/Popup";
import { useUserStore } from "@/store/user";
import JwtService from "@/services/JwtService";

export function useFollow() {
  const followLoading = ref(false);
  const isFollowing = ref(false);
  const followerCount = ref(0);
  const followers = ref<number[]>([])

  const fetchFollowerCount = async (userId: string | number) => {
    try {
      const [err, response] = await apiService.get<any>(`followers/${userId}`);
      if(err) return;
      followerCount.value = response.data.count;
      console.log('Follower count:', followerCount.value);
    } catch (error) {
      console.error("Takipçi sayısı alınırken hata oluştu:", error);
    }
  }

  const checkIfFollowing = async (userId: string | number) => {
    const userStore = useUserStore();
    // Only check if user is logged in
    if (!userStore.user.id) return;

    try {
      const [err, response] = await apiService.get<any>(`followers/${userId}`);
      if(err) return;

      const followers = response.data.followers || [];
      isFollowing.value = followers.some((follower: any) => follower.id === userStore.user.id);
      console.log('Is following:', isFollowing.value);
    } catch (error) {
      console.error("Takip durumu kontrol edilirken hata oluştu:", error);
    }
  }

  const toggleFollow = async (userId: string | number) => {
    const userStore = useUserStore();
    if (!userStore.user.id) {
      await router.push('/auth/login');
      return;
    }

    try {
      followLoading.value = true;

      if (isFollowing.value) {
        // Unfollow
        const [err] = await apiService.delete<any>(`follow/${userId}`);
        if(!err) {
          isFollowing.value = false;
          followerCount.value--;
        }
      } else {
        // Follow
        const [err] = await apiService.post<any>(`follow/${userId}`);
        if(!err) {
          isFollowing.value = true;
          followerCount.value++;
        }
      }
    } catch (error) {
      ErrorPopup('İşlem sırasında bir hata oluştu');
      console.error("Takip işlemi sırasında hata:", error);
    } finally {
      followLoading.value = false;
    }
  }

  const fetchFollowings = async (userId: string | number) => {
    try {
      console.log('Fetching followers for user ID:', userId)
      const userStore = useUserStore();
      
      // Kullanıcı giriş yapmış mı kontrol et
      if (!userStore.isAuthenticated) {
        console.warn('Kullanıcı giriş yapmamış, takip edilenler alınamıyor')
        followers.value = []
        return
      }
    
      
      const [err, response] = await apiService.get<any>(`followings/${userId}`)
      if (err) {
        console.error('Takip edilenler alınırken API hatası:', err)
        followers.value = []
        return
      }
      
      
      // Yeni API yanıt yapısı:
      // { following: number[], count: number, error_code: number, error_message: string }
      if (response && response.data && Array.isArray(response.data.following)) {
        // Doğrudan ID dizisi alalım
        followers.value = response.data.following
        console.log('Takip edilen kullanıcı ID\'leri:', followers.value)
      } else {
        console.warn('API beklenmeyen format döndü:', response)
        followers.value = []
      }
    } catch (error) {
      console.error("Takip edilenler alınırken hata oluştu:", error)
      followers.value = []
    }
  }

  return {
    followLoading,
    isFollowing,
    followerCount,
    fetchFollowerCount,
    checkIfFollowing,
    toggleFollow,
    fetchFollowings,
    followers
  };
}