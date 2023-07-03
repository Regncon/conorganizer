"use client";

import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Button, Card } from "../lib/mui";
import React, { Fragment, useContext, useEffect, useState } from "react";
import { faCoffee } from "@fortawesome/free-solid-svg-icons";
import {
  doc,
  onSnapshot,
  updateDoc,
  setDoc,
  deleteDoc,
  collection,
  serverTimestamp,
  getDocs,
  query,
  where,
  orderBy,
  limit,
} from "firebase/firestore";
import db from "../lib/firebase";
import { AuthContext } from "../components/auth";

interface Props {}

const Experiment = (props: Props) => {
  const colletionRef = collection(db, "schools");
  const [clicked, setClicked] = useState(false);

  const { currentUser } = useContext(AuthContext);

  const currentUserId = currentUser ? currentUser.uid : null;
  const [schools, setSchools] = useState([]);
  const [loading, setLoading] = useState(false);
  const [title, setTitle] = useState("");
  const [desc, setDesc] = useState("");
  const [score, setScore] = useState("");

  useEffect(() => {
    const q = query(
      colletionRef,
      //  where('owner', '==', currentUserId),
      where("title", "==", "School1") // does not need index
      //  where('score', '<=', 100) // needs index  https://firebase.google.com/docs/firestore/query-data/indexing?authuser=1&hl=en
      // orderBy('score', 'asc'), // be aware of limitations: https://firebase.google.com/docs/firestore/query-data/order-limit-data#limitations
      // limit(1)
    );

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

  // ADD FUNCTION
  const addSchool = async () => {
    const owner = currentUser ? currentUser.uid : "unknown";
    const ownerEmail = currentUser ? currentUser.email : "unknown";

    const newSchool = {
      title,
      desc,
      score: +score,
      owner,
      ownerEmail,
      createdAt: serverTimestamp(),
      lastUpdate: serverTimestamp(),
    };

    try {
      const schoolRef = doc(colletionRef);
      await setDoc(schoolRef, newSchool);
    } catch (error) {
      console.error(error);
    }
  };

  //DELETE FUNCTION
  async function deleteSchool(school) {
    try {
      const schoolRef = doc(colletionRef, school.id);
      await deleteDoc(schoolRef, schoolRef);
    } catch (error) {
      console.error(error);
    }
  }

  // EDIT FUNCTION
  async function editSchool(school) {
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
    <div className="flex flex-col items-center justify-center h-screen">
      <Button
        variant="contained"
        onClick={() => setClicked(!clicked)}
        className="flex items-center gap-2"
      >
        hello World
      </Button>
      {clicked && <FontAwesomeIcon icon={faCoffee} />}

      <Fragment>
        <h1>Schools (SNAPSHOT adv.)</h1>
        <div className="inputBox">
          <h3>Add New</h3>
          <h6>Title</h6>
          <input
            className="text-black"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
          <h6>Score 0-10</h6>
          <input
            className="text-black"
            type="number"
            value={score}
            onChange={(e) => setScore(e.target.value)}
          />
          <h6>Description</h6>
          <textarea
            className="text-black"
            value={desc}
            onChange={(e) => setDesc(e.target.value)}
          />
          <button onClick={() => addSchool()}>Submit</button>
        </div>
        <hr />
        {loading ? <h1>Loading...</h1> : null}
        {schools.map((school) => (
          <Card className="school card p-1 m-1" key={school.id}>
            <h2>{school.title}</h2>
            <p>{school.desc}</p>
            <p>{school.score}</p>
            <p>{school.ownerEmail}</p>
            <div>
              <button onClick={() => deleteSchool(school)}>X</button>
              <button onClick={() => editSchool(school)}>Edit Score</button>
            </div>
          </Card>
        ))}
      </Fragment>
    </div>
  );
};

export default Experiment;
