#!/usr/bin/env python3
"""
Import migration articles into the news table.

Usage:
    python3 import_news.py

Requires: psycopg2 (pip install psycopg2-binary)
Reads: /home/devbox/project/migration_output/articles.json
Writes to: PostgreSQL news table
"""

import json
import os
import re
from datetime import datetime

import psycopg2

ARTICLES_PATH = "/home/devbox/project/migration_output/articles.json"
IMAGE_BASE_URL = "/uploads/migration/images/"
MEDIA_BASE_URL = "/uploads/migration/media/"

DB_HOST = os.environ.get("DB_HOST", "xinhang-db-postgresql.ns-0h7fttt7.svc")
DB_PORT = os.environ.get("DB_PORT", "5432")
DB_USER = os.environ.get("DB_USER", "postgres")
DB_PASSWORD = os.environ.get("DB_PASSWORD", "7dl72vft")
DB_NAME = os.environ.get("DB_NAME", "xinhang")

CATEGORY_MAP = {
    "校园动态": "campus",
    "新闻动态": "news",
    "名师风采": "teachers",
    "学生活动": "activities",
    "学子风采": "highlights",
    "媒体聚焦": "media",
    "教育科研": "research",
    "教学教研": "research",
    "教学活动": "research",
    "科研动态": "research",
    "德育天地": "moral_education",
    "德育活动": "moral_education",
    "All": "news_en",
    "友好学校": "partner_schools",
    "家校交流": "home_school",
    "家校平台": "home_school",
    "家委会建设": "home_school",
    "父母大学": "home_school",
    "党建工作": "party",
    "招生专栏": "admission",
    "国际部": "international",
    "学生社团": "activities",
    "心灵花园": "moral_education",
}


def clean_html(html):
    """Clean inline styles from HTML content while preserving structure."""
    html = re.sub(r'\s*style="[^"]*"', '', html)
    html = re.sub(r'\s*style=\'[^\']*\'', '', html)
    html = html.replace('images/', IMAGE_BASE_URL)
    html = html.replace('media/', MEDIA_BASE_URL)
    return html


def extract_summary(html, max_len=200):
    """Extract plain text summary from HTML."""
    text = re.sub(r'<[^>]+>', '', html)
    text = re.sub(r'\s+', ' ', text).strip()
    text = text.replace('&nbsp;', ' ').replace('&amp;', '&')
    text = text.replace('&lt;', '<').replace('&gt;', '>')
    if len(text) > max_len:
        text = text[:max_len] + "..."
    return text


def parse_date(date_str):
    """Parse date string into datetime object."""
    for fmt in ("%Y-%m-%d %H:%M", "%Y-%m-%d %H:%M:%S", "%Y-%m-%d"):
        try:
            return datetime.strptime(date_str, fmt)
        except ValueError:
            continue
    return datetime.now()


def main():
    print("Loading articles.json...")
    with open(ARTICLES_PATH, encoding="utf-8") as f:
        articles = json.load(f)
    print(f"  Loaded {len(articles)} articles")

    # Deduplicate by title
    seen_titles = set()
    unique_articles = []
    for art in articles:
        title = art["title"].strip()
        if title in seen_titles:
            continue
        seen_titles.add(title)
        unique_articles.append(art)
    print(f"  After dedup: {len(unique_articles)} unique articles ({len(articles) - len(unique_articles)} duplicates removed)")

    # Connect to database
    print(f"Connecting to database {DB_NAME}@{DB_HOST}...")
    conn = psycopg2.connect(
        host=DB_HOST, port=DB_PORT,
        user=DB_USER, password=DB_PASSWORD,
        dbname=DB_NAME
    )
    cur = conn.cursor()

    # Check if news table has data already
    cur.execute("SELECT COUNT(*) FROM news")
    existing = cur.fetchone()[0]
    if existing > 0:
        print(f"  WARNING: news table already has {existing} rows")
        print("  Clearing existing imported data (old_id IS NOT NULL)...")
        cur.execute("DELETE FROM news WHERE old_id IS NOT NULL AND old_id != ''")
        conn.commit()
        print("  Cleared.")

    # Import
    imported = 0
    errors = 0
    for art in unique_articles:
        try:
            title = art["title"].strip()
            content = clean_html(art["content_html"])
            summary = extract_summary(art["content_html"])
            author = art.get("author", "").strip()
            source = art.get("source", "").strip()
            keywords = art.get("keywords", "").strip()
            old_id = art.get("article_id", "")
            column_name = art.get("column_name", "")
            category = CATEGORY_MAP.get(column_name, "other")

            cover = art.get("cover_image", "")
            if cover:
                cover = IMAGE_BASE_URL + cover.replace("images/", "")

            pub_time = parse_date(art.get("pub_date", ""))

            cur.execute("""
                INSERT INTO news (title, summary, content, category, cover_image,
                                  author_name, source, keywords, old_id,
                                  published, published_at, created_at, updated_at)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            """, (
                title, summary, content, category, cover,
                author, source, keywords, old_id,
                True, pub_time, pub_time, pub_time
            ))
            imported += 1
        except Exception as e:
            errors += 1
            print(f"  ERROR importing '{art.get('title', '?')[:30]}': {e}")

    conn.commit()
    cur.close()
    conn.close()

    print(f"\nDone!")
    print(f"  Imported: {imported}")
    print(f"  Errors: {errors}")
    print(f"  Categories used: {len(set(CATEGORY_MAP.get(a.get('column_name', ''), 'other') for a in unique_articles))}")


if __name__ == "__main__":
    main()
