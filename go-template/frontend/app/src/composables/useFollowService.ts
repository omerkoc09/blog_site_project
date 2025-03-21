import apiService from "@/services/ApiService";
import { ref } from "vue";
import { router } from "@/plugins/1.router";
import { ErrorPopup } from "@/utils/Popup";
import { useUserStore } from "@/store/user";

export function useFollow() {
  const followLoading = ref(false);
  const isFollowing = ref(false);
  const followerCount = ref(0);
  
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
      router.push('/auth/login');
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

  return {
    followLoading,
    isFollowing,
    followerCount,
    fetchFollowerCount,
    checkIfFollowing,
    toggleFollow
  };
}