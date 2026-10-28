import { Button, Form } from "react-bootstrap";

type PatternFilterProps = {
  dateFrom: string;
  dateTo: string;
  titleQuery: string;
  onDateFromChange: (value: string) => void;
  onDateToChange: (value: string) => void;
  onTitleQueryChange: (value: string) => void;
  onApply: () => void;
};

export function PatternFilter(props: PatternFilterProps) {
  return (
    <Form
      className="filter-bar"
      onSubmit={(event) => {
        event.preventDefault();
        props.onApply();
      }}
    >
      <Form.Group className="filter-bar__field">
        <Form.Label>С</Form.Label>
        <Form.Control
          type="date"
          value={props.dateFrom}
          onChange={(event) => props.onDateFromChange(event.target.value)}
        />
      </Form.Group>
      <Form.Group className="filter-bar__field">
        <Form.Label>по</Form.Label>
        <Form.Control
          type="date"
          value={props.dateTo}
          onChange={(event) => props.onDateToChange(event.target.value)}
        />
      </Form.Group>
      <Form.Group className="filter-bar__field filter-bar__field--title">
        <Form.Label>Название</Form.Label>
        <Form.Control
          type="text"
          placeholder="фрагмент названия"
          value={props.titleQuery}
          onChange={(event) => props.onTitleQueryChange(event.target.value)}
        />
      </Form.Group>
      <Button className="filter-bar__submit" type="submit">
        Найти
      </Button>
    </Form>
  );
}
