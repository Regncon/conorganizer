import Image from "next/image";
import styles from "./page.module.scss";
import { Box, Button } from "@mui/material";

export default function Home() {
  return (
    <Box component="main" className={styles.main}>
      Hello world! :D
    </Box>
  );
}
