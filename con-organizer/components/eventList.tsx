"use client";

import { Box, Button } from "../lib/mui";
import React, { useEffect, useState } from "react";
import {
  onSnapshot,
  collection,
} from "firebase/firestore";
import db from "../lib/firebase";
import EventUi from "./eventUi";

interface Props {}

const EventList = () => {
  const colletionRef = collection(db, "schools");
  const [schools, setSchools] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {

    setLoading(true);
    const unsub = onSnapshot(colletionRef, (querySnapshot) => {
      const items = [];
      querySnapshot.forEach((doc) => {
        items.push(doc.data());
      });
      setSchools(items);
      setLoading(false);
    });
    return () => {
      unsub();
    };
  }, []);

    // EDIT FUNCTION
    async function editEvent(school) {
        const updatedSchool = {
          score: +score,
          lastUpdate: serverTimestamp(),
        };
    
        try {
          const schoolRef = doc(colletionRef, school.id);
          updateDoc(schoolRef, updatedSchool);
        } catch (error) {
          console.error(error);
        }
      }
    

  return (
      <Box className="gap-4" >
        {loading ? <h1>Loading...</h1> : null}
        {schools.map((school) => (
            <EventUi title={school.title} image={school.image} description={school.description} />
        ))}
      </Box>
  );
};

export default EventList;
