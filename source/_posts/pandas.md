---
title: pandas
date: 2023-06-25 16:37:56
tags:
- Pandas
categories:
- Python
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---

## Summary

I learn pandas in kaggle recently, I feel that this is a simplified version of some commonly used SQL.


## Code

### Creating, Reading and Writing

```python
import pandas as pd

# create a simple DataFrame
pd.DataFrame({ "Yes": [50, 21], "No": [131, 2] })

# create a DataFrame with string value
pd.DataFrame({ "Bob": ["I liked it.", "It was awful"], "Sue": ["Pretty good.", "Bland"] })

# specify the index
pd.DataFrame({
    "Bob": [ "I liked it.", "It was awful."],
    "Sue": [ "Pretty good.", "Bland."],
}, index=["Product A", "Product B"])

# create a simple Series
pd.Series([ 1, 2, 3, 4, 5 ])

# specify the index
pd.Series([30, 35, 40], index=["2015 Sales", "2016 Sales", "2017 Sales"])

# read data from csv file
reviews = pd.read_csv("file_path")

# check how large the resulting DataFrame is
reviews.shape

# examine the contents of the resultant DataFrame
reviews.head()

# specify the index column when open the file
reviews = pd.read_csv("file_path", index_col=0)

# save the DataFrame to disk
reviews.to_csv("file_path")
```


### Indexing, Selecting and Assigning

```python
import pandas as pd

reviews = pd.read_csv("file_path", index_col=0)
pd.set_option("display.max_rows", 5)

# access the column `country`
reviews.country

# or
reviews["country"]

# select the first value in column `country`
reviews["country"][0]

# select the first row of data in a DataFrame
reviews.iloc[0]

# get a column with `iloc`
reviews.iloc[:, 0]

# select the `country` column(the first column) from just the first, second and third row
reviews.iloc[:4, 0]

# or, just select the second and third entries
reviews.iloc[1:3, 0]

# or pass a list
reviews.iloc[[1,2,3], 0]

# start count forwards from the end of the values
reviews.iloc[-5:]


# get the first entry in `reviews.country`
reviews.loc[0, "country"]

# select the data with specified columns
reviews.loc[:, ["taster_name", "taster_twitter_handle", "points"]]

# specify the index
reviews.set_index("title")

# conditional selection
# example, if the value if Italy of not
reviews.country == "Italy"

# select the relevant data
reviews.loc[reviews.country == "Italy"]

# more condition, and
reviews.loc[(reviews.country == "Italy") & (reviews.points >= 95)]

# or
reviews.loc[(reviews.country == "Italy") | (reviews.points >= 90)]

# a built-in conditional selector: isin
reviews.loc[reviews.country.isin(["Italy", "France"])]

# isnull, notnull
reviews.loc[reviews.price.notnull()]


# Assigning data
reviews["tmp"] = "everyone"

# with an iterable of values
reviews["index_backwards"] = range(len(reviews), 0, -1)
```


### Summary Functions and Maps

```python
import pandas as pd

df = pd.read_csv("file_path", index_col=0)

# summary functions, e.g. describe
df.points.describe()

# It's type-aware, its output changes based on the data type of the input

# mean
df.points.mean()

# to see a list of unique values: unique
df["taster_name"].unique()

# to see a list of unique values and how often they occur in the dataset
df["taster_name"].value_counts()


# Maps
# suppose that we wanted to remean the scores the wines received to 0
review_points_mean = df.points.mean()
df.points.map(lambda p: p - review_points_mean)

# `map` processes Series by single value
# `apply` processes DataFrame by row
# Not that `map()` and `apply()` return new, transformed Series and DataFrame
# They don't modify the original data they're called on.

def remean_points(row):
    row.points = row.points - review_points_mean
    return row
    
df.apply(remean_points, axis="columns")


# All of the standard Python operators (`>, <, ==`, and so on) work in this manner
review_points_mean = df.points.mean()
df.points - review_points_mean
```


### Grouping and Sorting

```python
import pandas as pd

df = pd.read_csv("file_path", index_col=0)

# groupwise analysis
df.groupby("points").points.count()
# equal: df.points.value_counts()

# get the minimum value in each point value category
df.groupby("points").price.min()

# get the best wine with points by country and province
df.groupby(["country", "province"]).apply(lambda df: df.loc[df.points.idxmax()])

# agg: run a bunch of different functions
# generate a simple statistical summary of the dataset
df.groupby(["country"]).price.agg([len, min, max])

# multi-indexes
country_review = df.groupby(["country", "province"]).description.agg([len])

# convert back to a regular index
country_review.reset_index()


# sorting
country_review.sort_values(by="len")

# `sort_values` defaults to an ascendint sort, if we want a descending sort
country_review.sort_values(by="len", ascending=False)

# sort by index values
country_review.sort_index()

# sort by more than one column at a time
country_review.sort_values(by=["country", "len"])
```


### Data Types and Missing Values

```python
import pandas as pd

df = pd.read_csv("file_path", index_col=0)
pd.set_option("display.max_rows", 5)

# get the dtype of the `price` column in the df
df.price.dtype

# get the dtype of every column in the DataFrame
df.dtypes

# convert a column of one type into another wherever such conversion makes sense
df.points.astype("float64")

# a DataFrame or Series index has its own dtype
df.index.dtype

# select the NaN entries
df[pd.isnull(df.country)]

# replacing missing values with `fillna`
df.region_2.fillna("Unknown")

# replace a non-null value
# e.g. replace UA to UK
df["country".replace("UA", "UK")]
```


### Renaming and Combining

```python
import pandas as pd

df = pd.read_csv("file_path", index_col=0)

# change index names and/or column names
# e.g. change the points column to score
df.rename(columns={"points": "score"})

# change the index
df.rename(index={0: "firstEntry", 1: "secondEntry"})

# Both the row index and the column index can have their own `name` attribute
# use `rename_axis()` to change these names
df.rename_axis("wines", axis="rows").rename_axis("fields", axis="columns")


# combine data from different place: concat, join, merge

dfA = pd.read_csv("file_a_path")
dfB = pd.read_csv("file_b_path")

pd.concat([dfA, dfB])

# 
left = dfA.set_index(["title", "trending_date"])
right = dfB.set_index(["title", "trending_date"])

left.join(right, lsuffix="_A", rsuffix="_B")

```
