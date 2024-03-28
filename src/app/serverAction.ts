'use server';
import { db } from "$lib/firebase/firebase";
import { adminDb } from "$lib/firebase/firebaseAdmin";
import { doc, getDocs } from "firebase/firestore";
import { EventCardProps } from "./EventCardProps";

export async function getByID(id: string) {
    const eventRef = adminDb.collection('event').doc(id);
const doc = await eventRef.get();
if (!doc.exists) {
  console.log('No such document!');
  return null;
} else {
  console.log('Document data:', doc.data());
  return doc.data();
}
}
export async function getAll() {
  const eventRef = await adminDb.collection('event').get();
  return eventRef.docs.map(doc => doc.data()) as Omit<EventCardProps, 'icons'>[];
}