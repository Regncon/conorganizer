import AppTopBar from "@/components/AppTopBar";
import EventUi from "@/components/eventUi";
import AddEvent from "@/components/addEvent";

export default function Home() {
  return (
    <main className="">
      <AppTopBar />
      <EventUi />
      <EventUi />
      <EventUi />
      {/*       <div>
        <AuthProvider>
          <Login />
          <Welcome />
          <Experiment />
        </AuthProvider>
      </div> */}
      <AddEvent />
    </main>
  );
}
