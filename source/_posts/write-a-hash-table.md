---
title: Write a hash table
date: 2023-06-21 13:39:31
tags:
- C
- Golang
categories:
- Tutorial
- Data Structure
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---

## Summary

This is a tiny little demo of hash table. But it can already be used to understand the map in golang in some aspects.


## Steps

- define data structure
- implement hash function
- API
  - insert
  - search
  - delete


## Code

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define TABLE_SIZE 100
#define LOAD_FACTOR_THRESHOLD 0.75

typedef struct HashNode
{
  char *key;
  char *value;
  struct HashNode *next;
} HashNode;

typedef struct HashTable
{
  int size;
  int count;
  HashNode **buckets;
} HashTable;

void rehash(HashTable *table);

unsigned int hash(char *key)
{
  unsigned int hash = 0;
  while (*key)
  {
    hash = (hash << 5) + *key++;
  }
  return hash % TABLE_SIZE;
}

HashTable *create_table()
{
  HashTable *table = malloc(sizeof(HashTable));
  table->size = TABLE_SIZE;
  table->count = 0;
  table->buckets = calloc((size_t)table->size, sizeof(HashNode*));
  return table;
}

void del_table(HashTable* ht)
{
  free(ht->buckets);
  free(ht);
}

void insert(HashTable *table, char *key, char *value)
{
  unsigned int idx = hash(key);
  HashNode *newNode = malloc(sizeof(HashNode));
  newNode->key = strdup(key);
  newNode->value = strdup(value);
  newNode->next = table->buckets[idx];
  //
  table->buckets[idx] = newNode;
  table->count++;
  //
  float loadFactor = (float)table->count / table->size;

  if (loadFactor > LOAD_FACTOR_THRESHOLD)
  {
    rehash(table);
  }
}

char *search(HashTable *table, char *key)
{
  unsigned int idx = hash(key);
  HashNode *node = table->buckets[idx];
  while (node)
  {
    if (strcmp(node->key, key) == 0)
    {
      return node->value;
    }
    node = node->next;
  }
  return NULL;
}

void delete(HashTable *table, char *key)
{
  unsigned int idx = hash(key);
  HashNode *node = table->buckets[idx];
  HashNode *prev = NULL;
  while (node)
  {
    if (strcmp(node->key, key) == 0)
    {
      if (prev)
      {
        prev->next = node->next;
      }
      else
      {
        table->buckets[idx] = node->next;
      }
      free(node->key);
      free(node->value);
      free(node);
      table->count--;
      return;
    }
    prev = node;
    node = node->next;
  }
}

void rehash(HashTable *table)
{
  int old_size = table->size;
  HashNode **oldBuckets = table->buckets;

  table->size *= 2;
  table->buckets = calloc((size_t)table->size, sizeof(HashNode *));

  for (int i = 0; i < old_size; i++)
  {
    HashNode *node = oldBuckets[i];
    while (node)
    {
      insert(table, node->key, node->value); // 使用新的哈希表大小重新插入
      HashNode *oldNode = node;
      node = node->next;
      free(oldNode->key);
      free(oldNode->value);
      free(oldNode);
    }
  }

  free(oldBuckets);
}


int main()
{
  HashTable *table = create_table();

  insert(table, "name", "Alice");
  insert(table, "age", "25");

  printf("Name: %s\n", search(table, "name"));
  printf("Age: %s\n", search(table, "age"));
  printf("ht->size: %d\tht->count: %d\n", table->size, table->count);

  delete(table, "name");

  printf("Name after deletion: %s\n", search(table, "name"));
  printf("ht->size: %d\tht->count: %d\n", table->size, table->count);
  del_table(table);
  return 0;
}
```
