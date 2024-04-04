import { useState } from "react";
import { v4 as uuidv4 } from "uuid";

const url = "ws://localhost:8080/connect";
const ws = new WebSocket(url);
let myid = "";

const sizeOfSide = 10;

class FieldsValues {
  static empty = 0;
  static placed = 1;
  static missed = 2;
  static hit = 3;
}

interface IntroMessage {
  type: "intro_message";
  number_of_connected_players: number;
}

interface FieldMessage {
  type: "change_field";
  first_field: number[];
  second_field: number[];
}

interface ActionMessage {
  change_to: number
  changed_positions: number[]
}

interface Cell {
  column: number;
  row: number;
}

function cellPosToNumber(cell: Cell) {
  return cell.column * 10 + cell.row
}

const strings = {
  phases: {
    setup: "Setup phase",
  },
  notification: {
    playersNotConnected: {
      first: "Waiting for the first player",
      second: "Waiting for the second player",
    },
    playersConnected: {
      first: "First player connected",
      second: "Second player connected",
    },
  },
};

export default function App() {
  const [myField, setMyField] = useState<number[][]>();
  const [enemyField, setEnemyField] = useState<number[][]>();
  const [currentPhase, setCurrentPhase] = useState<string>(
    strings.phases.setup
  );
  const [numberOfConnectedPlayers, setNumberOfConnectedPlayers] =
    useState<number>(0);

  function handleClick(cell: Cell) {
    const m: ActionMessage = {
      change_to: FieldsValues.hit,
      changed_positions: [cellPosToNumber(cell)]
    }
    ws.send(JSON.stringify(m));
  }

  ws.addEventListener("message", (event: MessageEvent<string>) => {
    const json: IntroMessage | FieldMessage = JSON.parse(event.data);
    if (json.type === "intro_message") {
      console.log("Recieved message: ", json);
      setNumberOfConnectedPlayers(json.number_of_connected_players);
      if (json.number_of_connected_players === 2) {
        setMyField(getBlankField())
        setEnemyField(getBlankField())
      }
    } else if (json.type === "change_field") {
      console.log(json.first_field, oneDimensionToTwoDimensions(json.first_field, sizeOfSide))
      setMyField(oneDimensionToTwoDimensions(json.first_field, sizeOfSide))
      setEnemyField(oneDimensionToTwoDimensions(json.second_field, sizeOfSide))
    }
  });

  return (
    <div>
      <h1 className="phase">{currentPhase}</h1>
      <div className="wrapper">
        <div>
          {myField ? (
            myField.map((arr, indexTop) => {
              return (
                <div className="row">
                  {arr.map((val, indexLow) => {
                    const cell: Cell = {
                      column: indexTop,
                      row: indexLow,
                    };
                    return (
                      <div
                        className="cell"
                        id={uuidv4()}
                        onClick={() => handleClick(cell)}
                      >
                        {val}
                      </div>
                    );
                  })}
                </div>
              );
            })
          ) : (
            <p className="notification">
              {numberOfConnectedPlayers === 1 || numberOfConnectedPlayers === 2
                ? strings.notification.playersConnected.first
                : strings.notification.playersNotConnected.first}
            </p>
          )}
        </div>
        <div>
          {enemyField ? (
            enemyField.map((arr, _) => {
              return (
                <div className="row">
                  {arr.map((val, _) => {
                    return (
                      <div
                        className="cell"
                        id={uuidv4()}
                      >
                        {val}
                      </div>
                    );
                  })}
                </div>
              );
            })
          ) : (
            <p className="notification">
              {numberOfConnectedPlayers === 2
                ? strings.notification.playersConnected.second
                : strings.notification.playersNotConnected.second}
            </p>
          )}
        </div>
      </div>
    </div>
  );
}

function oneDimensionToTwoDimensions(arr: number[], size: number) {
  let arr1: number[] = [];
  const arr2: number[][] = [];
  for (const i in arr) {
    arr1.push(arr[i]);
    if (+i % size === size - 1) {
      arr2.push(arr1);
      arr1 = [];
    }
  }
  return arr2;
}

function getBlankField() {
  return Array(sizeOfSide)
    .fill(FieldsValues.empty)
    .map(() => Array(sizeOfSide).fill(FieldsValues.empty));
}
