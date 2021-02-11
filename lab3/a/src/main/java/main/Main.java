package main;

import java.util.ArrayDeque;
import java.util.concurrent.*;
import java.util.concurrent.atomic.*;
import java.util.concurrent.locks.ReentrantLock;


class Pot {
    public Integer value = 0;
    public Semaphore sync = new Semaphore(0);
}

class Bear extends Thread {
    volatile Pot pot;

    public Bear(Pot p) {
        pot = p;
    }

    public synchronized void run() {
            while (true) {
                try {
                    wait();

                    pot.value = 0;

                    System.out.println("Bear eats");
                    Thread.sleep(100);

                    pot.sync.release();

                } catch (InterruptedException err) {
                    System.err.println(err);
                }
            }
    }
}

class Bee extends Thread {
    volatile Pot pot;
    volatile Bear bear;
    static volatile ReentrantLock sync = new ReentrantLock();

    public Bee(Pot p, Bear b) {
        pot = p;
        bear = b;
    }

    public void run() {
        while (true) {
            try {
                sync.lock();

                pot.value += 1;
                System.out.println("Bee " + Thread.currentThread().getName() + " " + pot.value);

                if (pot.value == 10) {

                    synchronized(bear) {
                        bear.notify();
                    }

                    pot.sync.acquire();
                }

                sync.unlock();

                Thread.sleep(1000);

            } catch (InterruptedException err) {
                System.err.println(err);
            }
        }
    }
}

public class Main {
    public static void main(String[] args) {
        int N = 3;

        Pot pot = new Pot();
        Bear bear = new Bear(pot);

        bear.start();

        for (int i = 0; i < N; i++) {
            Bee b = new Bee(pot, bear);
            b.start();
        }
    }
}