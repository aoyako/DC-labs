package main;

import java.util.ArrayDeque;
import java.util.concurrent.*;
import java.util.concurrent.atomic.*;
import java.util.concurrent.locks.*;
import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.sql.Time;
import java.util.Random;


class RWLock implements Lock {
    private int readers = 0;
    private int writers = 0;

    public RWLock() {
        
    }

    public synchronized void read() {
        while (writers != 0) {
            try {
                wait();
            } catch (InterruptedException err) {
                System.err.println(err);
            }
        }
        notifyAll();
    }

    public synchronized void write() {
        while ((readers != 0) && (writers != 0)) {
            try {
                wait();
            } catch (InterruptedException err) {
                System.err.println(err);
            }
        }
        notify();
    }

    public synchronized void unlock() {
        if (readers != 0) {
            readers--;
        } else if (writers != 0) {
            writers--;
        }
    }
}

class FilePrinter extends Thread {
    int[][] field;
    RWLock lock;

    FilePrinter(int[][] f, RWLock l) {
        field = f;
        lock = l;
    }

    public void run() {
        try {
            FileWriter file = new FileWriter("log.txt");
            while (true) {
                lock.read();

                StringBuilder b = new StringBuilder();
                for (int i = 0; i < field.length; i++) {
                    for (int j = 0; j < field[0].length; j++) {
                        b.append(field[i][j]);
                    }
                    b.append("\n");
                }
                b.append("\n");
                file.write(b.toString());
                file.flush();
                lock.unlock();

                Thread.sleep(1000);
            }

        } catch (Exception err) {
            System.err.println(err);
        }
    }
}

class ConsolePrinter extends Thread {
    int[][] field;
    int n = 0;
    RWLock lock;

    ConsolePrinter(int[][] f, RWLock l) {
        field = f;
        lock = l;
    }

    public void run() {
        try {
            while (true) {
                lock.read();

                System.out.println(n);
                StringBuilder b = new StringBuilder();
                for (int i = 0; i < field.length; i++) {
                    for (int j = 0; j < field[0].length; j++) {
                        b.append(field[i][j]);
                    }
                    b.append("\n");
                }

                System.out.println(b.toString());
                lock.unlock();
                n++;

                Thread.sleep(1000);
            }

        } catch (Exception err) {
            System.err.println(err);
        }
    }
}

class Gardener extends Thread {
    int[][] field;
    RWLock lock;

    Gardener(int[][] f, RWLock l) {
        field = f;
        lock = l;
    }

    public void run() {
        try {
            while (true) {
                Thread.sleep(3000);

                lock.read();

                for (int i = 0; i < field.length; i++) {
                    for (int j = 0; j < field[0].length; j++) {
                        if (field[i][j] > 0) {
                            lock.unlock();
                            lock.write();

                            field[i][j] = 0;

                            lock.unlock();
                            continue;
                        }
                    }
                }

                lock.unlock();
            }

        } catch (Exception err) {
            System.err.println(err);
        }
    }
}

class Nature extends Thread {
    int[][] field;
    RWLock lock;

    Nature(int[][] f, RWLock l) {
        field = f;
        lock = l;
    }

    public void run() {
        try {
            Random rnd = new Random();
            while (true) {
                lock.write();

                int x = rnd.nextInt(7);
                int y = rnd.nextInt(7);

                field[x][y] += 5;
                field[x+1][y] += 5;
                field[x][y+1] += 5;
                field[x+1][y+1] += 5;

                lock.unlock();

                Thread.sleep(2000);
            }

        } catch (Exception err) {
            System.err.println(err);
        }
    }
}

public class Main {
    public static void main(String[] args) {
        int[][] field = new int[10][10];
        RWLock lock = new RWLock();

        FilePrinter fp = new FilePrinter(field, lock);
        ConsolePrinter cp = new ConsolePrinter(field, lock);
        Gardener g = new Gardener(field, lock);
        Nature n = new Nature(field, lock);

        fp.start();
        cp.start();
        g.start();
        n.start();
    }
}