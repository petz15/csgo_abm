"""
Batch analysis script for tournament matchup folders.
Automatically analyzes all matchup folders and generates HTML reports.
"""

import os
import subprocess
import json
import argparse
from pathlib import Path
from datetime import datetime
from concurrent.futures import ProcessPoolExecutor, as_completed
import sys

def find_matchup_folders(tournament_folder):
    """Find all matchup_XXX folders in the tournament directory."""
    matchup_folders = []
    tournament_path = Path(tournament_folder)
    
    if not tournament_path.exists():
        print(f"Error: Tournament folder '{tournament_folder}' not found!")
        return []
    
    for item in sorted(tournament_path.iterdir()):
        if item.is_dir() and item.name.startswith('matchup_'):
            csv_file = item / "all_games_minimal.csv"
            json_file = item / "simulation_summary.json"
            
            if csv_file.exists() and json_file.exists():
                matchup_folders.append(item)
            else:
                print(f"⚠️  Skipping {item.name} - missing CSV or JSON file")
    
    return matchup_folders

def analyze_matchup(matchup_folder, notebook_path):
    """
    Analyze a single matchup by executing the notebook with updated configuration.
    Returns (success, matchup_name, report_path, error_msg)
    """
    matchup_name = matchup_folder.name
    
    try:
        # Paths
        csv_path = matchup_folder / "all_games_minimal.csv"
        json_path = matchup_folder / "simulation_summary.json"
        
        # Read strategy names from JSON
        with open(json_path, 'r') as f:
            sim_data = json.load(f)
        
        t1_strat = sim_data['simulation_config']['team1_strategy']
        t2_strat = sim_data['simulation_config']['team2_strategy']
        
        # Generate output filename
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        output_html = matchup_folder / f"analysis_report_{timestamp}_{t1_strat}_vs_{t2_strat}.html"
        
        # Create a temporary notebook with updated paths using papermill-style parameters
        # We'll use nbconvert with preprocessing
        
        # Execute notebook using jupyter nbconvert (optimized for speed)
        cmd = [
            'jupyter', 'nbconvert',
            '--to', 'html',
            '--no-input',  # Hide code cells
            '--execute',
            '--ExecutePreprocessor.timeout=300',  # 5 minute timeout (reduced from 10)
            '--ExecutePreprocessor.kernel_name=python3',  # Explicit kernel
            '--output', str(output_html),
            str(notebook_path)
        ]
        
        # Set environment variables for the notebook to read
        env = os.environ.copy()
        env['MATCHUP_FOLDER'] = str(matchup_folder)
        env['CSV_FILE'] = str(csv_path)
        env['JSON_FILE'] = str(json_path)
        env['MPLBACKEND'] = 'Agg'  # Use non-interactive matplotlib backend (faster)
        
        # Execute
        result = subprocess.run(
            cmd,
            env=env,
            capture_output=True,
            text=True,
            cwd=matchup_folder.parent  # Run from parent directory
        )
        
        if result.returncode == 0:
            return (True, matchup_name, str(output_html), None)
        else:
            error_msg = result.stderr[-500:] if result.stderr else "Unknown error"
            return (False, matchup_name, None, error_msg)
            
    except Exception as e:
        return (False, matchup_name, None, str(e))

def main():
    # Parse command line arguments
    parser = argparse.ArgumentParser(
        description='Batch analyze tournament matchup folders and generate HTML reports.',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog='''
Examples:
  python batch_analyze_tournament.py results_20260105_123550
  python batch_analyze_tournament.py results_20260105_123550 --workers 8
  python batch_analyze_tournament.py results_20260105_123550 --notebook analysis_notebook.ipynb
''')
    
    parser.add_argument('tournament_folder', 
                        help='Path to the tournament results folder containing matchup_XXX subdirectories')
    parser.add_argument('--notebook', '-n',
                        default='analysis_notebook_v2.ipynb',
                        help='Path to the Jupyter notebook to execute (default: analysis_notebook_v2.ipynb)')
    parser.add_argument('--workers', '-w',
                        type=int,
                        default=4,
                        help='Number of parallel workers (default: 4)')
    
    args = parser.parse_args()
    
    # Configuration from arguments
    TOURNAMENT_FOLDER = args.tournament_folder
    NOTEBOOK_PATH = Path(args.notebook).resolve()  # Convert to absolute path
    MAX_WORKERS = args.workers
    
    # Validate inputs
    if not Path(TOURNAMENT_FOLDER).exists():
        print(f"❌ Error: Tournament folder '{TOURNAMENT_FOLDER}' does not exist!")
        return 1
    
    if not NOTEBOOK_PATH.exists():
        print(f"❌ Error: Notebook '{NOTEBOOK_PATH}' does not exist!")
        return 1
    
    print("=" * 80)
    print("TOURNAMENT BATCH ANALYSIS (Parallel)")
    print("=" * 80)
    print(f"Tournament Folder: {TOURNAMENT_FOLDER}")
    print(f"Notebook: {NOTEBOOK_PATH}")
    print(f"Parallel Workers: {MAX_WORKERS}")
    print()
    
    # Find all matchup folders
    matchup_folders = find_matchup_folders(TOURNAMENT_FOLDER)
    
    if not matchup_folders:
        print("❌ No valid matchup folders found!")
        return 1
    
    print(f"✓ Found {len(matchup_folders)} matchup folders to analyze\n")
    
    # Track results
    successful = []
    failed = []
    start_time = datetime.now()
    
    # Process in parallel
    with ProcessPoolExecutor(max_workers=MAX_WORKERS) as executor:
        # Submit all jobs
        futures = {
            executor.submit(analyze_matchup, folder, NOTEBOOK_PATH): folder 
            for folder in matchup_folders
        }
        
        # Process as they complete
        for i, future in enumerate(as_completed(futures), 1):
            folder = futures[future]
            try:
                success, name, report_path, error = future.result()
                
                if success:
                    successful.append((name, report_path))
                    print(f"✓ [{i}/{len(matchup_folders)}] {name} - SUCCESS")
                else:
                    failed.append((name, error))
                    print(f"✗ [{i}/{len(matchup_folders)}] {name} - FAILED")
                    if error:
                        print(f"    Error: {error[:100]}")
            except Exception as e:
                failed.append((folder.name, str(e)))
                print(f"✗ [{i}/{len(matchup_folders)}] {folder.name} - EXCEPTION: {e}")
    
    # Summary
    elapsed = datetime.now() - start_time
    print("\n" + "=" * 80)
    print("BATCH ANALYSIS COMPLETE")
    print("=" * 80)
    print(f"Total Time: {elapsed}")
    print(f"Successful: {len(successful)}/{len(matchup_folders)}")
    print(f"Failed: {len(failed)}/{len(matchup_folders)}")
    
    if successful:
        print(f"\n✓ Successfully analyzed {len(successful)} matchups")
    
    if failed:
        print(f"\n✗ Failed matchups:")
        for name, error in failed[:10]:  # Show first 10 failures
            print(f"  - {name}: {error[:100]}")
        if len(failed) > 10:
            print(f"  ... and {len(failed) - 10} more")
    
   
    
    return 0 if not failed else 2  # Exit code 2 if some analyses failed



if __name__ == "__main__":
    sys.exit(main())
