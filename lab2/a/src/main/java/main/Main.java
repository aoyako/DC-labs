package main;

import java.util.ArrayDeque;
import java.util.concurrent.*;
import java.util.concurrent.atomic.*;

class MySyncronized {

    // AtomicInteger indicator = new AtomicInteger(0);
    // Object waiter = new Object();
    Semaphore sync = new Semaphore(1);

    public void lock() {
        try {
            sync.acquire();
        } catch (InterruptedException err) {
            System.err.println(err);
        }
    }

    public void unlock() {
        sync.release();
    }
}

class QueenBee {
    private boolean found = false;
    ArrayDeque<Integer> tasks = new ArrayDeque<Integer>();
    MySyncronized sync = new MySyncronized();

    public QueenBee(int fieldSize, int taskSize) {
        for (int i = 0; i < fieldSize; i += taskSize) {
            tasks.addLast(i);
        }
    }

    public Integer GetNextTask() {
        sync.lock();

        if (found) {
            sync.unlock();
            return null;
        }

        Integer val = tasks.pollFirst();

        sync.unlock();

        return val;
    }

    public void Found() {
        sync.lock();

        found = true;

        sync.unlock();
    }
}

class Bee extends Thread {
    volatile QueenBee queen;
    String name;
    int[][] field;
    int taskSize;

    public Bee(QueenBee q, int[][] field, int taskSize, String name) {
        this.queen = q;
        this.field = field;
        this.name = name;
        this.taskSize = taskSize;
    }

    public void run() {
        while (true) {
            Integer pos = queen.GetNextTask();
            if (pos != null) {
                System.out.println(name + " got task: " + pos);
                for (int i = pos; i < pos + taskSize; i++) {
                    for (int j = 0; j < field.length; j++) {
                        if (field[i][j] == 1) {
                            System.out.println(name + " found bear at " + String.valueOf(i) + "," + String.valueOf(j));
                            queen.Found();
                            System.out.println(name + " finishes");
                            return;
                        }
                    }
                }
            } else {
                System.out.println(name + " finishes");
                return;
            }
        }
    }
}

public class Main {

    static int[][] createField(int size, int bearX, int bearY) {
        int[][] field = new int[size][size];

        field[bearX][bearY] = 1;

        return field;
    }
    public static void main(String[] args) {
        int fieldSize = Integer.parseInt(args[0]);
        int taskSize = Integer.parseInt(args[1]);
        int beeNumber = Integer.parseInt(args[2]);

        int[][] field = createField(fieldSize, Integer.parseInt(args[3]), Integer.parseInt(args[4]));
        QueenBee queen = new QueenBee(fieldSize, taskSize);

        Bee[] bees = new Bee[beeNumber];
        for (int i = 0; i < beeNumber; i++) {
            bees[i] = new Bee(queen, field, taskSize, "Bee-" + String.valueOf(i));
            bees[i].start();
        }

        for (int i = 0; i < beeNumber; i++) {
            try {
                bees[i].join();
            } catch (InterruptedException err) {
                System.err.println(err);
            }
        }
    }
}