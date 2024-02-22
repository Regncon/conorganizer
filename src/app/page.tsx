import Image from "next/image";
import styles from "./page.module.scss";
import { Box, Button, CssBaseline, ThemeProvider } from "@mui/material";
import { muiDark } from "@/lib/muiTheme";

export default function Home() {
  return (
			<Box component="main" className={styles.main}>
				Hello world! :D
			</Box>
  );
}
