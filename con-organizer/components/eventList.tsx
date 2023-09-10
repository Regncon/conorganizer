"use client";

import { Box, Card, CardContent, CardHeader, CardMedia, Chip } from "../lib/mui";
import React, { useEffect, useState } from "react";
import {
  onSnapshot,
  collection,
} from "firebase/firestore";
import db from "../lib/firebase";
import EventUi from "./eventUi";
import { ConEvent } from "@/lib/types";
import { faDiceD20 } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

interface Props { }

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
  }, [conEvents]);

  return (
    <Box className="flex flex-row flex-wrap justify-center gap-4">
      {loading ? <h1>Loading...</h1> : null}
      {conEvents.map((conEvent) => (
        <>
          {/* <EventUi key={conEvent.id} colletionRef={colletionRef} conEvent={conEvent} /> */}
          <Card key={conEvent.id}>
            <CardHeader
              sx="padding-bottom: 0.5rem;"
              title={conEvent?.title}
              subheader="Kjempebra spennende event."
            />
            <Box className="flex justify-start pb-4">
            <CardMedia
                className="ml-4"
                sx="width: 40%; max-height: 130px"
                component="img"
                image="/placeholder.jpg"
                alt={conEvent?.title}
              />
              <Box
              className="flex flex-col pl-4 pr-4" >
                <span> 
                  
                <Chip icon={<FontAwesomeIcon icon={faDiceD20} />}
                  label="Rollespill"
                  size="small"
                  variant="outlined"
                  />
                  </span>
                <span>DnD 5e </span>                <span>
                Rom 222,
                </span>
                <span>
                  Søndag: 12:00 - 16:00
                </span>
              </Box>
            </Box>
          </Card>
        </>
      ))}
    </Box>
  );
};

export default EventList;
