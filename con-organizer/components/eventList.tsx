"use client";

import { Box } from "../lib/mui";
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

  return (
      <Box>
        {loading ? <h1>Loading...</h1> : null}
        {schools.map((school) => (
            <EventUi title={school.title} image={school.image} description={school.description} />
        ))}
      </Box>
  );
};

export default EventList;
