'use client';

// Import the functions you need from the SDKs you need
import { initializeApp } from 'firebase/app';
import { getAuth } from 'firebase/auth';
//import { getAnalytics } from "firebase/analytics";
import { getFirestore } from 'firebase/firestore';

// TODO: Add SDKs for Firebase products that you want to use
// https://firebase.google.com/docs/web/setup#available-libraries

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
const firebaseConfig = {
    apiKey: 'AIzaSyCcBmX5Sn-Tru0zd-5LUKFs-fuNhMGg7uw',
    authDomain: 'regncon2023.firebaseapp.com',
    projectId: 'regncon2023',
    storageBucket: 'regncon2023.appspot.com',
    messagingSenderId: '667813245424',
    appId: '1:667813245424:web:84c9f9d360368ac089ed35',
    measurementId: 'G-K8EKFKXE9N',
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
//const analytics = getAnalytics(app);

export const auth = getAuth();

export default getFirestore(app);
