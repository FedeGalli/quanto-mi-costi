// stores/auth.js
import { writable } from "svelte/store";
import { onAuthStateChanged, signOut } from "firebase/auth";
import { auth } from "./credentials";

// Create the user store
export const user = writable(null);
export const isAuthenticated = writable(false);
export const isLoading = writable(true);

// Initialize auth state listener
let unsubscribe: any;

export function initAuthStore() {
  if (unsubscribe) return; // Prevent multiple listeners

  unsubscribe = onAuthStateChanged(auth, (firebaseUser) => {
    if (firebaseUser && firebaseUser.emailVerified) {
      // User is signed in and email is verified
      const userData: any = {
        uid: firebaseUser.uid,
        email: firebaseUser.email,
        displayName: firebaseUser.displayName,
        photoURL: firebaseUser.photoURL,
        emailVerified: firebaseUser.emailVerified,
        firstName: firebaseUser.displayName?.split(" ")[0] || "",
        lastName: firebaseUser.displayName?.split(" ").slice(1).join(" ") || "",
      };

      user.set(userData);
      isAuthenticated.set(true);
    } else {
      // User is signed out or email not verified
      user.set(null);
      isAuthenticated.set(false);
    }
    isLoading.set(false);
  });
}

// Logout function
export async function logout() {
  try {
    await signOut(auth);
    user.set(null);
    isAuthenticated.set(false);
  } catch (error) {
    console.error("Logout error:", error);
    throw error;
  }
}

// Cleanup function
export function destroyAuthStore() {
  if (unsubscribe) {
    unsubscribe();
    unsubscribe = null;
  }
}
