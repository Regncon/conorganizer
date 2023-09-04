"use client";

import { Box, Button } from "../lib/mui";
import React, { useEffect, useState } from "react";
import {
  onSnapshot,
  collection,
} from "firebase/firestore";
import db from "../lib/firebase";
import EventUi from "./eventUi";
import { ConEvent } from "@/lib/types";

interface Props {}

const EventList = () => {
  const colletionRef = collection(db, "schools");
  const [conEvents, setconEvents] = useState([] as ConEvent[]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {

    setLoading(true);
    const unsub = onSnapshot(colletionRef, (querySnapshot) => {
      const items = [] as ConEvent[];
      querySnapshot.forEach((doc) => {
        items.push(doc.data());
      });
      setconEvents(items);
      setLoading(false);
    });
    return () => {
      unsub();
    };
  }, []); 

  useEffect(() => {
    console.log("conEvents", conEvents);
  } , [conEvents]);

  return (
      <Box className="gap-4" >
        {loading ? <h1>Loading...</h1> : null}
        {conEvents.map((conEvent) => (
            <EventUi key={conEvent.id} colletionRef={colletionRef} conEvent={conEvent} />
        ))}
      </Box>
  );
};

export default EventList;
