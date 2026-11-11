import { Spinner } from "react-bootstrap";
import { useAppSelector } from "../store/hooks";

export function BusyOverlay() {
  const authLoading = useAppSelector((state) => state.auth.loading);
  const patternsLoading = useAppSelector((state) => state.patterns.loading);
  if (!authLoading && !patternsLoading) {
    return null;
  }
  return (
    <div className="busy" role="status" aria-live="polite">
      <Spinner animation="border" variant="light" />
    </div>
  );
}
