"use client";

import { Box } from "../lib/mui";
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
        items.push(doc.data() as ConEvent);
        items[items.length - 1].id = doc.id;
      });
      setconEvents(items);
      setLoading(false);
    });
    return () => {
      unsub();
    };
  }, []); 

  useEffect(() => {
  } , [conEvents]);

  return (
      <Box className="flex flex-row flex-wrap justify-center gap-4">
        {loading ? <h1>Loading...</h1> : null}
        {conEvents.map((conEvent) => (
            <EventUi key={conEvent.id} colletionRef={colletionRef} conEvent={conEvent} />
            <h1>Hei Gerhard, dette er testen til Christer, du fant den!</h1>
        ))}
      </Box>
  );
};

export default EventList;
